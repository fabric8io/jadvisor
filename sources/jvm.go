package sources

var (
	JVMRequest = JolokiaRequest{
		Type:      Read,
		MBean:     "java.lang:type=Memory",
		Attribute: []string{"HeapMemoryUsage", "NonHeapMemoryUsage"},
	}
)
