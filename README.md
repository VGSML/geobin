# GeoBin - Working with Bitmap Indexed Spatial Data in Go

GeoBin is a GoLang package that provides a simple and efficient way to work with spatial data sets (such as GeoJSON) using bitmap indexes.

## Key Features

- **Integration with Uber's H3Geo**: Utilizes [uber/h3geo](https://h3geo.org) for aggregation, intersection, contains, and joining operations with geometry (spatial attributes).
- **Geometry Data Structure**: Employs [github.com/paulmach/orb](https://github.com/paulmach/orb) for managing geometry data structures.
- **Bitmap Indexing**: Implements bitmap indexing using [github.com/roaring/roaring](https://github.com/roaring/roaring).

## Subpackages

- **bjoin**: Manages the data structure to store results of joining two bitmap indexed data sets.
- **h3b**: Bitmap index specifically for H3 cells.
- **h3f**: Functions for working with H3 cells in Go, extending the uber/h3 (v4) package ([https://github.com/uber/h3-go](https://github.com/uber/h3-go)).
- **orbf**: Offers functions and operations for geometries.

## Testing and Fixtures

For testing purposes, GeoBin uses data from OpenStreetMap ([https://openstreetmap.org](https://openstreetmap.org)) to ensure reliability and accuracy.

## Contributing

We welcome contributions to GeoBin and appreciate your efforts to improve the project. Here’s how you can contribute:

### Reporting Issues

- **Bug Reports**: If you encounter a bug, please open an issue with a detailed description and a code sample or an executable test case demonstrating the unexpected behavior.
- **Feature Requests**: Feel free to suggest new features by opening an issue with a thorough explanation of the feature and its potential benefits.

### Making Changes

- **Fork the Repository**: Begin by forking the GeoBin repository on GitHub.
- **Create a Branch**: Make a new branch in your fork for your changes.
- **Implement Your Changes**: Implement your changes in your branch. Make sure to follow the existing coding style. Add or update tests as necessary.
- **Test Your Changes**: Ensure that your changes do not break any existing functionality and that all tests pass.
- **Submit a Pull Request**: Once you're ready, submit a pull request to the main GeoBin repository with a clear description of the problem and your solution.

### Code Style and Conventions

- We follow the [Uber Go Style Guide](https://github.com/uber-go/guide). Please ensure your contributions adhere to these guidelines for consistency and readability.

### Code of Conduct

- All contributors are expected to adhere to our Code of Conduct. Please be respectful and considerate in your interactions within the GeoBin community.

Thank you for contributing to GeoBin!

## License

GeoBin is made available under the MIT License. This license allows you to use, share, modify, and distribute the package and any derivative works, provided that you give appropriate credit to the original author.

### MIT License

Copyright (c) 2024 Gribanov Vladimir

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.