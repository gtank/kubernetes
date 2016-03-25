package certificate

import "testing"

const (
	testCSR = "-----BEGIN CERTIFICATE REQUEST-----\nMIIDLDCCAhQCAQAwgaMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlh\nMRYwFAYDVQQHDA1TYW4gRnJhbmNpc2NvMSkwJwYDVQQKDCBIb25lc3QgQWNobWVk\ncyBVc2VkIENlcnRpZmljYXRlczEmMCQGA1UECwwdSGFzdGlseSBHZW5lcmF0ZWQg\nVmFsdWVzIERlcHQxFDASBgNVBAMMC2t1YmUtd29ya2VyMIIBIjANBgkqhkiG9w0B\nAQEFAAOCAQ8AMIIBCgKCAQEAtQHE+2NQW2UXU0W9bd4MvJJbv+Hm4HhUqP6X18Pm\n06bmb00KpVDQbakA5HzxSDibMA5KlPo0OwF912VKGaYLDdzX4hCsnF2+zgiZXBAH\nn4NVQDCyPJb35G4u6aV789wn5t28RPYBso2yPNBs/2oXwwCm38HpasFmkGLrtp4S\nRqBDyNmxTMY9qtYZKHl/ImbD2101XZmTUqFJlK1+3M0gPgZPeGTLCGBHOEOxSDWH\nYAq1MGtvSiym05hEjlgWstUbon84EX3lXROz+qS/9nyewiL6PXJHJWMyObS2tkQg\n+ipJkhXKrdNkBzQ4arhZZaaKjjYpQedim4J3MQ356Jz1CwIDAQABoEMwQQYJKoZI\nhvcNAQkHMTQMMlNIQTI1Njo2U1FKbTRSMktjRXJDUDFaa0libUo5QmdNM3pNTXpx\nejkrYjdrVHMxRi9rMA0GCSqGSIb3DQEBCwUAA4IBAQCZU1WxJS17aN3BFWWfwNHg\nBMA+7UmMM2WKPpkvOA7YAW4MfsZDktwLej+4bkiuVkFHsrYsA+vHHTBEyRhOikPU\nfi/x0zKuf4Sb5DhS74OWfErghwI14rLOQvnO5Fzzzq42odn9izsLqOn8cALa6K93\n3Gj1QoRWWTsoeB46PZfZuDusivymghkVactJOedkTwQNDygNl2GHIiV2fIycjoN5\nJuN2WQvrmMEpjF4PQtic8jEL56D9bpxsg5ljAPDqS5XZWqQ5yGK8hSC7QSRWPEuc\n9HCrTasyuO8/mmnvxiS8HQ/e1cnTTPVKNeviQX3Oqdsm5m47YdqGCx9MX+8Ok/60\n-----END CERTIFICATE REQUEST-----"
)

func TestDecodeCertificateRequest(t *testing.T) {
	csr, err := DecodeCertificateRequest([]byte(testCSR))
	if err != nil {
		t.Fatal(err)
	}
	if csr.Subject.CommonName != "kube-worker" {
		t.Error("read wrong CN from certificate request")
	}
}
