FROM centos:7

COPY ./k8s-controller-custom-resource /
RUN chmod +x /k8s-controller-custom-resource

WORKDIR /
CMD ["/k8s-controller-custom-resource"]
