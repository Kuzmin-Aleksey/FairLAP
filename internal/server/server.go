package server

type Server struct {
	detector  *DetectorServer
	groups    *GroupsServer
	metrics   *MetricServer
	lapConfig *LapConfigServer
	images    *ImagesServer
}

func NewServer(
	detector *DetectorServer,
	groups *GroupsServer,
	metrics *MetricServer,
	lapConfig *LapConfigServer,
	images *ImagesServer,
) *Server {
	return &Server{
		detector:  detector,
		groups:    groups,
		metrics:   metrics,
		lapConfig: lapConfig,
		images:    images,
	}
}
