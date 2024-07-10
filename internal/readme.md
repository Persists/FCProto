This is the folder for the internal shared packages

- `models`: Contains the message model which is used to serialize and deserialize messages. It contains other useful models.
- `queue`: Contains the logic for the message queue which is used to buffer messages before they are sent and after they are received.
- `connection`: Contains the logic for the connection between the cloud and edge application. It utilizes the message queue to send and receive messages. It also uses tcp to establish a connection between the cloud and edge application. It uses the message model to serialize and deserialize messages.
- `utils`: Contains utility functions that are used in the cloud and edge application.
