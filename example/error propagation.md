# Error Progapation

With concurrent code, and especially distributed systems, it's both easy for something to go wrong in your system, and difficult to understand why it happened.

Many developers make the mistake of thinking of error propagation as secondary, or "other" to the flow of their system.

Errors indicate that your system has entered a state in which it cannot fulfill an operation that a user either explicitly requested. Because of this, it needs to relay a few pieces of critical information:

**What happened**
This is the part of the error that contains information about what happened, e.g., "disk full", "socket closed" or "credentials expired"

**When and where it occurred**

Error should always contain a complete stack trace starting with how the call was initiated and ending with where the error was instantiated. The stack trace should not be caontained in the error message (more on this in a bit), but should be easily accessible when handlilng the error up the stack.

Further, the error should contain information regarding the context it's running within. For example, in a distributed system, it should have some way of identifying what machine the error occurred on. Later, when trying to understand what happend in your system, this information will be invaluable.

In addition, the error should contain the time on the machine the error was instantiated on, in UTC.

**A friendly user-facing message**
The message that gets displayed to the user should be customized to suit your system and its users.

**How the user can get more information**
At some point, someone will likely want to know, in detail, what happened when the error occurred. Errors that are presented to users should provide an ID that can be cross-referenced to a corresponding log that displays the full information of the error.

It's possible to place all errors into one of two categories:

- Bugs
- Known edge cases (e.g., broken network connections, failed disk writes, etc.)