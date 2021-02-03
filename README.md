# Pluralsight

This repository follows [this](https://app.pluralsight.com/library/courses/go-build-distributed-applications/table-of-contents) course on pluralsight.

# Startup

1. Start the Website in `src\distributed\web`.
2. Start the Coordinator in `src\distributed\coordinator\exec`.
3. Start the Sensors using the `start_sensors.bat` script.

# Description

Each sensor creates its own data queue and publishes a message to a discovery queue.
Also, sensors reply to Discover Requests with their queue name.

The queue listener listens for messages on the discovery fanout exchange. For each sensor,
a Goroutine is created. Those goroutines react to messages on the sensor queues by publishing
an event via the Event Aggregator.

The DB and Webapp consumer are subscribed to those events and can now use that data. 
The DB consumer for example listens to the MessageReceived_sensor_name event after
having discovered the sensor queue. The event data is then put into the PersistData-Queue.

Now the (seperate instance!) Datamanager can subscribe to the PersistData-Queue and write
the event data to the database.
