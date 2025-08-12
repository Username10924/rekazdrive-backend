package storage

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

type S3Backend struct {
	Endpoint string // endpoint URL
	Bucket   string // S3 bucket name
	AccessKey string
	Secret string
	Region string
	client *http.Client // HTTP client for making requests
}

func NewS3Backend(endpoint, bucket, accessKey, secret, region string) *S3Backend {
	return &S3Backend{
		Endpoint: strings.TrimRight(endpoint, "/"),
		Bucket:   bucket,
		AccessKey: accessKey,
		Secret:   secret,
		Region:   region,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *S3Backend) objectURL(object string) string {
	// <endpoint>/<bucket>/<object>
	return s.Endpoint + "/" + s.Bucket + "/" + path.Clean("/"+object)[1:]
}

func (s *S3Backend) Save(id string, data []byte) error {
	obj := sanitizeID(id)
	urlStr := s.objectURL(obj)

	req, err := http.NewRequest("PUT", urlStr, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Host", req.URL.Host)
	req.Header.Set("Content-Type", "application/octet-stream") // binary data stream

	if err := s.signV4(req, data, time.Now().UTC()); err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to save object %s: status: %d body: %s", obj, resp.StatusCode, string(b))
	}

	return nil // success
}

func (s *S3Backend) Load(id string) ([]byte, error) {
	obj := sanitizeID(id)
	urlStr := s.objectURL(obj)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", req.URL.Host)

	if err := s.signV4(req, []byte{}, time.Now().UTC()); err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to load object %s: status: %d body: %s", obj, resp.StatusCode, string(b))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil // success
}

func (s *S3Backend) Delete(id string) error {
	obj := sanitizeID(id)
	urlStr := s.objectURL(obj)

	req, err := http.NewRequest("DELETE", urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Host", req.URL.Host)

	if err := s.signV4(req, []byte{}, time.Now().UTC()); err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete object %s: status: %d body: %s", obj, resp.StatusCode, string(b))
	}
	
	return nil // success
}


/* ------------ AWS Signature V4 implementation ------------ */
func (s *S3Backend) signV4(req *http.Request, payload []byte, t time.Time) error {
	service := "s3"

	amzDate := t.UTC().Format("20060102T150405Z")
	dateStamp := t.UTC().Format("20060102")

	// payload hash
	payloadHash := sha256.Sum256(payload)
	payloadHashHex := hex.EncodeToString(payloadHash[:])
	req.Header.Set("x-amz-content-sha256", payloadHashHex)
	req.Header.Set("x-amz-date", amzDate)

	// canonical URL
	canonicalURL := req.URL.Path
	// canonical query
	var keys []string
	for k := range req.URL.Query() {
		keys = append(keys, k)
	}
	sort.Strings(keys) // sorting is a must in aws signature v4, param1=value1&param2=value2 is different from param2=value2&param1=value1
	canonicalQuery := ""
	for i, k := range keys {
		if i > 0 {
			canonicalQuery += "&"
		}
		canonicalQuery += url.QueryEscape(k) + "=" + url.QueryEscape(req.URL.Query().Get(k))
	}

	// canonical headers
	canonicalHeaders := "host:" + req.URL.Host + "\n" +
		"x-amz-content-sha256:" + payloadHashHex + "\n" +
		"x-amz-date:" + amzDate + "\n"
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"

	// canonical request
	canonicalRequest := req.Method + "\n" +
		canonicalURL + "\n" +
		canonicalQuery + "\n" +
		canonicalHeaders + "\n" +
		signedHeaders + "\n" +
		payloadHashHex
	
	hash := sha256.Sum256([]byte(canonicalRequest))
	stringsToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, s.Region, service),
		hex.EncodeToString(hash[:]),
	}, "\n")

	signingKey := s.getSigningKey(dateStamp, s.Region, service)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringsToSign))

	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, s.Region, service)
	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s", s.AccessKey, credentialScope, signedHeaders, signature)

	req.Header.Set("Authorization", authHeader)
	return nil
}

// helper function that hashes input data using SHA256, returns hex-encoded string
func hmacSHA256(key []byte, data string) []byte {
	m := hmac.New(sha256.New, key)
	m.Write([]byte(data))
	return m.Sum(nil)
}

// helper function that generates the signing key for AWS Signature V4
func (s *S3Backend) getSigningKey(dateStamp, region, service string) []byte {
	// AWS Signature Version 4 signing key derivation
	// 1. Create a date key using the secret and the date
	// 2. Create a region key using the date key and the region
	// 3. Create a service key using the region key and the service
	// 4. Create a signing key using the service key and "aws4_request"
	kDate := hmacSHA256([]byte("AWS4"+s.Secret), dateStamp) // valid only for the date
	kRegion := hmacSHA256(kDate, region) // valid only for the date and region
	kService := hmacSHA256(kRegion, service) // valid only for the date, region, and service
	kSigning := hmacSHA256(kService, "aws4_request") // final signing key (reusable for all requests with same date/region/service)
	return kSigning
}
