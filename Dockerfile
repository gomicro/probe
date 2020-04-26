FROM scratch
MAINTAINER dev@gomicro.io

ADD probe probe

CMD ["/probe"]
