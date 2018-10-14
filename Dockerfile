FROM golang:1.11

RUN apt update -y

RUN apt install -y python

RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py

RUN python get-pip.py

RUN pip install tensorflow numpy

WORKDIR /go/src/app

COPY pyscripts pyscripts

RUN /usr/bin/python pyscripts/classify_image.py

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 3001

CMD ["app"]
