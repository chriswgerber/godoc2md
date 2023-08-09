PKG_NAME:=godoc2md
PROJECT_NAME=github.com/chriswgerber/godoc2md

EXE=./$(PKG_NAME)

TARGETS=$(PKG_NAME)

$(PKG_NAME): $(wildcard *.go)
	go build -o $@ ./cmd/$@

build: $(PKG_NAME)

all: readme examples doc

readme: README.md

doc: README.md

examples:
	$(EXE) github.com/kr/fs > examples/fs/README.md
	$(EXE) github.com/codegangsta/martini > examples/martini/README.md
	$(EXE) github.com/gorilla/sessions > examples/sessions/README.md
	$(EXE) go/build > examples/build/README.md

.PHONY: examples readme all

README.md: $(PKG_NAME)
	$(EXE) -v github.com/chriswgerber/godoc2md > $@
