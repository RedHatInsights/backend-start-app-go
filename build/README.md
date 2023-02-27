# `/build`

Building, Packaging and Continuous Integration info.

## Packaging

Keep your container (Docker), OS (deb, rpm, pkg) package configurations and scripts in the `/build` directory.

## CI

Put your CI (travis, circle, drone) configurations and scripts in the `/build/ci` directory.
Note, that some of the CI tools (e.g., Travis CI) are very picky about the location of their config files. Try putting the config files in the `/build/ci` directory linking them to the location where the CI tools expect them when possible (don't worry if it's not and if keeping those files in the root directory makes your life easier :-)).
