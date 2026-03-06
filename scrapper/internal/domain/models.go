package domain

type Internship struct {
	CompanyName string
	SourceURL   string
	Tracks      []Track
}

type Track struct {
	PositionName string
	TechStack    []string
	MinSalary    int32
	Location     string
}
