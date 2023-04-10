package k8sutil

const (
	BorealisDNSNameAnnotation = "external-dns.alpha.kubernetes.io/hostname"
	ElbTimeoutAnnotationName  = "service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout"
	ElbTimeoutAnnotationValue = "3600"
)
