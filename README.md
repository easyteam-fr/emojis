# Emoji

Emoji is a sample application that works with a Kubernetes Operator to discover
custom resources that describes Emojis and displays them in a web page.

## Design

The application design looks like below:

![Design](img/emoji-design.png)

- An Emoji CRD allow to declare Emojis in Kubernetes
- A controller get all the Emojis from Kubernetes and synchronize them with the
  application. Once done, it publishes the status back to the resource.

## Building the application


