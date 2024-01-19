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
 * **matrix**: https://matrix.to/#/#image-builder:fedoraproject.org
*  **Mailing List**: image-builder@redhat.com

#### Contributing

Please refer to the [developer guide](https://www.osbuild.org/guides/developer-guide/index.html) to learn about our workflow, code style and more.

### License:

 - **Apache-2.0**
 - See LICENSE file for details.
