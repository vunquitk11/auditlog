package kafka

// getErrorForKafkaRetry standardizes the error returned from Kafka operations
// to determine if a retry should be attempted. Extend this function to handle
// specific error types that require custom retry logic.
func getErrorForKafkaRetry(err error) error {
	switch err {
	default:
		return err
	}
}
