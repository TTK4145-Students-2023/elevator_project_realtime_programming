Elevator Project
================
This is a program that controls an elevator and lets it communicate with other elevators to collaborate on completing orders. The program contains a database, which holds information about all the elevators in the system. Every elevator uses UDP broadcast to share the information about itself to the other elevators, which is updated to the database in each program. This lets all the elevators know everything about each other, hence they have the same world view.

An order is received by the buttons being pushed, and information about it is stored in the database. To decide which elevator who should conduct the order, a fleeting master is used. The elvator that receives the button press is the fleeting master, and it delegates the order. This is possible because the fleeting master know everything about the other elevators, and can therefore decide which elevator is best suited to take the order.

The system works dynamically, so that if one elevator malfunctions, the other elevators can take the orders that the malfuctioning elevator was meant to do. If the malfunctioning elevator starts working again, the other elevators will send which cab calls the elevator had before it malfunctioned.  


Here is a brief description of what the functionality of the modules are.

-SingleElevator: Controls elevator and does actions based on given orders. Controls physical aspects of the elevator.

-DatabaseHandler: Updating database based on messages received from other elevators. Assigner determines which elevator who should take orders, based on the info in database.

-ConnectionHandler: To see which elevators are connected to the system, performs actions when elevators disconnects and reconnects.

-Network: Establishes network connection and functions for sending messages across network.

-ElevatorHardware: Interface between physical elevator and program. Registering button presses, setting motor direction and likewise. 











