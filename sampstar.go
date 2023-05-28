package sampstar

type SampstarOptions struct {
	Port int
}

type SampstarServer struct {
	Options *SampstarOptions
}

func NewSampstarServer(options *SampstarOptions) *SampstarServer {
	return &SampstarServer{
		Options: options,
	}
}

func (sampstarServer *SampstarServer) Start() {}
