Community envoy gateway
=======================

Gateway for the image builder community service.

### Run the example

To run envoy:
```
docker run --net=host -v $PWD/example:/app -it envoyproxy/envoy:distroless-v1.29-latest -c /app/config.yaml
```

To try it out:
```
AT=$(cat example/auth/access-token)
curl -H "authorization: Bearer $AT" localhost:10000/
```

### Project

 * **Website**: <https://www.osbuild.org>
 * **Bug Tracker**: <https://github.com/osbuild/community-gateway/issues>
 * **Discussions**: <https://github.com/orgs/osbuild/discussions>
 * **Matrix**: [#image-builder on fedoraproject.org](https://matrix.to/#/#image-builder:fedoraproject.org)
 * **Changelog**: <https://github.com/osbuild/community-gateway/releases>

#### Contributing

Please refer to the [developer guide](https://osbuild.org/docs/developer-guide/index) to learn about our workflow, code style and more.

### License:

 - **Apache-2.0**
 - See LICENSE file for details.
