# FCProto
This repository contains a prototype developed for the Fog Computing lecture. The project demonstrates the key principles and practical applications of fog computing, including distributed data processing and edge device coordination.





### Components
The projects holds the following components:
- **Queue**: A simple queue implementation that can be used to store and retrieve data. The Idea behind this component is that when i message should be sent to the cloud, it is stored in the queue and then sent to the cloud when the connection is available. This way, the message is not lost when the connection is down. When the cloud receives the message it is put in the queue of the cloud and then processed. This way, the cloud can also store messages when the connection to the edge device is down. So both the edge device and the cloud have two queues, one for sending and one for receiving messages. Furthermore this queue should enable parallel processing of messages.
