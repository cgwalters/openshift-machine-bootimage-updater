FROM registry.svc.ci.openshift.org/openshift/release:golang-1.15 AS builder
WORKDIR /src
COPY . .
RUN make

FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
COPY --from=builder /src/openshift-update-bootimages /usr/bin
CMD ["/usr/bin/openshift-update-bootimages"]
