package world

type Materialmap []Material

type Material struct {
	Id    int
	Name  string
	Solid int
}

func LoadMaterialmap(file string) (*Materialmap, error) {
	mmap, err := unmarshall(file)
	if err != nil {
		return nil, err
	}
	materials, err := getChildrenArray(mmap, "materials")
	if err != nil {
		return nil, err
	}
	ret := Materialmap{}
	for _, mat := range materials {
		cur, err := mapToMaterial(mat)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *cur)
	}
	return &ret, nil
}

func mapToMaterial(data map[string]interface{}) (*Material, error) {
	id, err := getIntValue(data, "id")
	if err != nil {
		return nil, err
	}
	name, err := getStringValue(data, "name")
	if err != nil {
		return nil, err
	}
	solid, err := getIntValue(data, "solid")
	if err != nil {
		return nil, err
	}
	return &Material{Id: id, Name: name,
		Solid: solid}, nil
}
