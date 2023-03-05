$version: "2"
namespace example.weather

service Weather {
    version: "2006-03-01"
    resources: [City]
}

resource City {
    read: GetCity
    list: ListCities
}

operation GetCity {
    input: GetCityInput
    output: GetCityOutput
    errors: []
}

structure GetCityInput {
}

structure GetCityOutput {
    name: String
}
