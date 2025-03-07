# vulcan-frame

**vulcan-frame** is a game server framework that is ready to use. It provides a general server framework for microservices, allowing you to quickly build a high-performance and highly available game server cluster.

## vulcan-pkg-app

**vulcan-pkg-app** is a package that provides a general function for **vulcan-frame**. It provides the components of routing distribution and load balancing, link tracing, metrics, etc., which can be easily combined into **vulcan-frame**. It is based on [kratos](https://github.com/go-kratos/kratos).

## Components

- **route**: grpc connection with route distribution and load balancing
- **trace**: distributed tracing
- **metrics**: metrics collection
- **logger**: common logger
- **context**: common context
- **profile**: profile definition
- **errors**: common errors definition
- **version**: version generation tool

## Contributing

If you have any suggestions or feedback, please feel free to open an issue or submit a pull request.
