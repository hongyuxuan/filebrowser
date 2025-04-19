FROM alpine:3.19.1
ADD bin/filebrowser /opt
WORKDIR /opt
RUN chmod +x /opt/filebrowser
EXPOSE 9121
ENTRYPOINT [ "./filebrowser" ]
