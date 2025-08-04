FROM redhat/ubi9:9.4

ARG BPP_PATH

ENV BPP_PATH=${BPP_PATH}

WORKDIR /app

COPY ${BPP_PATH} ./bpp

RUN chmod +x ./bpp

EXPOSE 8080

CMD ["./bpp"]
