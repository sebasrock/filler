#!/bin/bash -eu
# **************************************************************
# ** Coverage Golang app                                      **
# ** Date : 10/07/2019                                        **
# ** Author: ADL DevOps Team                                  **
# **************************************************************
export LC_CTYPE=C

RED='\033[0;31m'
NC='\033[0m' # No Color

# Code coverage generation
COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
PKG_LIST=$(go list ./... | grep -v /vendor/)

coverageByPackage() {
  echo 'mode: count' >"${COVERAGE_DIR}"/coverage.cov
  # Create a coverage file for each package
  printf "${RED}Creating coverage file for each package...${NC}\n"
  for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package"
    if [[ -f ${COVERAGE_DIR}/"${package##*/}.cov" ]]; then
      tail -n +2 ${COVERAGE_DIR}/"${package##*/}.cov" >>${COVERAGE_DIR}/coverage.cov
    fi
  done
}

coverageTerminal() {
  # Display the global code coverage
  go tool cover -func="${COVERAGE_DIR}"/coverage.cov
}

coverageHtml() {
  # If needed, generate HTML report
  printf "${RED}Build coverage html...${NC}\n"
  go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o coverage.html
}

clean() {
  # Remove the coverage files directory
  printf "${RED}Clean coverage directory...${NC}\n"
  rm -rf "$COVERAGE_DIR"
}

main() {
  # Create the coverage files directory
  printf "${RED}Creating coverage directory...${NC}\n"
  mkdir -p "$COVERAGE_DIR"
  coverageByPackage
  coverageTerminal
  coverageHtml
  clean
  exit 0
}

main "$@"
