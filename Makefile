a:
	go run main.go gen -o ./a.sql --config ./rbacgen.yaml

b:
	go run main.go gen -o ./b.sql --config ./rbacgen.rest.yaml

c:
	go run main.go gen -o ./c.sql --config ./rbacgen.retailbase.yaml

all:
	make a
	make b
	make c