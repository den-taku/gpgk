# compiler
CXX = g++
CXXFLAGS = -O2 -Wall -std=c++11

# src directories
DIR_CLIENT = ./src
DIR_GRPC = ./grpc
DIR_TEST = ./test

# target directories
TARGET = target
TARGET_DEPS = ${TARGET}/deps

# src file
SRC := ${shell find . -depth 1 -name "*.cpp"}
SRC_CLIENT := ${shell find ${DIR_CLIENT} -depth 1 -name "*.cpp"}
SRC_GRPC := ${shell find ${DIR_GRPC} -depth 1 -name "*.cpp"}
SRC_TEST := ${shell find ${DIR_TEST} -depth 1 -name "*.cpp"}

SRCS = ${SRC} ${SRC_CLIENT} ${SRC_GRPC} ${SRC_TEST}
SRCS_LIB = ${SRC} ${SRC_GRPC}

# objects files
OBJS = ${foreach src,${SRCS},${patsubst %.cpp,${TARGET_DEPS}/%.o,${notdir ${src}}}}
OBJS_LIB = ${foreach src,${SRCS_LIB},${patsubst %.cpp,${TARGET_DEPS}/%.o,${notdir ${src}}}}
OBJS_TEST = ${foreach src,${SRC_TEST},${patsubst %.cpp,${TARGET_DEPS}/%.o,${notdir ${src}}}}
OBJS_CLIENT = ${foreach src,${SRC_CLIENT},${patsubst %.cpp,${TARGET_DEPS}/%.o,${notdir ${src}}}}

.PHONY: clean test run

run: ${TARGET}/client
	./$<

build: ${TARGET}/client

# run test
test: ${TARGET}/test
	./$<

${TARGET}/client: ${OBJS_LIB} ${OBJS_CLIENT}
	${CXX} ${CXXFLAGS} -o $@ $^

# make ELF for test
${TARGET}/test: ${OBJS_LIB} ${OBJS_TEST}
	${CXX} ${CXXFLAGS} -o $@ $^

# compile each src file to object file when file updated
define template
${TARGET_DEPS}/${patsubst %.cpp,%.o,${notdir $1}}: $1 ${TARGET_DEPS}
	if [ ! -e ${TARGET_DEPS}/${patsubst %.cpp,%.o,${notdir $1}} ] \
	|| [ ${TARGET_DEPS}/${patsubst %.cpp,%.o,${notdir $1}} -ot ${1} ]; then \
		${CXX} -c ${CXXFLAGS} -o ${TARGET_DEPS}/${patsubst %.cpp,%.o,${notdir $1}} ${1}; \
	fi
endef

${foreach src,${SRCS},${eval ${call template,${src}}}}

# remove binaries
clean:
	rm -rf ${TARGET}

# make target directory
${TARGET_DEPS}:
	mkdir -p $@