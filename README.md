# storage-orchestrator

## Package Overview

The `pkg` package provides an approach to storage orchestration, enabling complex operations for saving, retrieving, and deleting data across multiple storage units with customizable strategies and concurrency control.

## Types

### Options Types

Options for each method (Save, Get, Delete) allow for customization of storage operations. Each option function modifies an internal state that dictates how the operation should be performed.

#### SaveOptionsFunc

A function that modifies the save options. It uses a pointer to SaveOptions, allowing adjustments like:

- *Context*: The operation's context, used for cancellation and metadata propagation.
- *HowWillItSave*: Determines whether the operation will be Sequential or Parallel, affecting performance and execution order.
- *Targets*: Specifies the storage units to be used.

#### GetOptionsFunc

A function that modifies the get options. Allows adjustments such as:

- *Context*: Similar to SaveOptions.
- *HowWillItGet*: Defines the retrieval strategy (Cache for quick access or Race to wait for pending save operations).
- *Targets*: Specifies the specific units to be queried.

#### DeleteOptionsFunc

A function that adjusts the deletion options. Allows specifying:

- *Context*: For managing concurrent operations and cancellation.
- *HowWillItDelete*: Can be SequentialDelete, ensuring the execution order of deletions.
- *Targets*: Specific units where deletion should occur.



## Order of Operations

Move this section after the "Types" to flow logically from the detailed options types to how these options are utilized to set operation orders.

### Through Option Functions

Explain how the order can be set dynamically through option functions for each type of operation.

### Using SetStandardOrder Method

Discuss the method to set a default, global order, explaining its utility and application within the orchestrator.

## Risks and Considerations

### Error Handling

- Implement retry logic.
- Communicate errors clearly.
- Recommend returning specific errors when requested items or units are not found.

### Data Consistency

- Challenges in `Parallel` operations.
- Strategies for synchronization and rollback.

## Best Practices

- **Use Contexts**: How and why to use `context.Context`.
- **Logging and Monitoring**: Importance and methods for effective logging and monitoring.

## Example Usage

Direct users to the `examples` directory for practical illustrations and scenarios. Explain that this section is designed to enhance understanding through real-world application examples.

Hello Jão
