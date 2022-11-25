package constant

import "time"

const ProcessTimeout = 2 * time.Minute
const APITimeout = 10 * time.Second

const CoordinatorServiceURL = "http://coordinator:8080/api"