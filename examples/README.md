# SimBrief SDK Examples

This directory contains comprehensive examples demonstrating how to use the SimBrief SDK for various flight planning scenarios. Each example focuses on different aspects of the SDK and progresses from basic usage to advanced integration patterns.

## Available Examples

### 📚 [Basic Example](basic/)

**Perfect for getting started** - Demonstrates fundamental SDK usage including:

- Getting supported aircraft types and layouts
- Fetching existing flight plans by user ID
- Basic error handling and response parsing
- Simple configuration with environment variables

**Prerequisites**: Internet connection (SimBrief user ID optional)

```bash
cd basic
go run main.go
```

### 🚀 [Advanced Example](advanced/)

**Production-ready patterns** - Shows sophisticated features including:

- Custom aircraft configurations with detailed specifications
- Comprehensive flight plan building with all parameters
- Detailed flight plan analysis and data extraction
- Weather integration and performance calculations
- Weight & balance analysis and fuel planning
- File generation and download capabilities

**Prerequisites**: Internet connection + **Required** SimBrief user ID

```bash
cd advanced
go run main.go
```

## Quick Start Guide

### 1. Choose Your Starting Point

- **New to SimBrief?** → Start with the [Basic Example](basic/)
- **Ready for production?** → Jump to the [Advanced Example](advanced/)
- **Need specific features?** → See feature comparison below

### 2. Set Up Environment

```bash
# Optional: For fetching existing flight plans
export SIMBRIEF_USER_ID=857341

# Optional: Enable debug logging
export SIMBRIEF_DEBUG=true
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run Examples

```bash
# Basic example
cd basic && go run main.go

# Advanced example  
cd advanced && go run main.go
```

## Feature Comparison

| Feature | Basic | Advanced | Description |
|---------|-------|----------|-------------|
| **Getting Started** | ✅ | ✅ | Simple SDK initialization |
| **Aircraft Types** | ✅ | ✅ | Fetch supported aircraft |
| **Flight Plan Fetch** | ✅ | ✅ | Get existing flight plans |
| **Custom Aircraft** | ❌ | ✅ | Define custom aircraft specs |
| **Flight Plan Builder** | ❌ | ✅ | Comprehensive plan creation |
| **Weight & Balance** | ❌ | ✅ | Detailed W&B analysis |
| **Fuel Planning** | ❌ | ✅ | Complete fuel calculations |
| **Weather Data** | ❌ | ✅ | METAR and weather analysis |
| **Navigation Log** | ❌ | ✅ | Route and waypoint data |
| **File Downloads** | ❌ | ✅ | PDF, XML, flight sim files |
| **Error Handling** | Basic | ✅ | Production-ready patterns |

## Learning Path

### 1. Start with Basic (15 minutes)
- Understand SDK initialization
- Learn basic API calls
- See simple error handling
- Get familiar with data structures

### 2. Move to Advanced (30 minutes)
- Master complex request building
- Learn data extraction patterns
- Understand production patterns
- Explore integration possibilities

### 3. Build Your Own
- Combine patterns from both examples
- Add your specific requirements
- Implement custom features
- Deploy to production

## Common Use Cases

### Flight Operations Centers
```go
// Use advanced example patterns for:
// - Real-time flight planning
// - Performance analysis
// - Fuel optimization
// - Route planning
```

### Training Applications
```go
// Use basic example for:
// - Learning flight planning
// - Understanding aircraft types
// - Simple demonstrations
```

### Integration Projects
```go
// Combine both examples for:
// - Custom flight planning tools
// - Analytics platforms
// - Mobile applications
// - Web services
```

## Environment Setup

### Required Environment Variables

```bash
# For advanced example (required)
export SIMBRIEF_USER_ID=your_user_id_here

# Optional: Enable detailed logging
export SIMBRIEF_DEBUG=true
```

### Finding Your SimBrief User ID

1. Go to [SimBrief.com](https://simbrief.com)
2. Log in to your account
3. Navigate to Account Settings
4. Find your User ID in the account information section

### Alternative: Using the Web Interface

If you don't have a user ID, you can still:
- Use the basic example to explore aircraft types
- Generate flight plan URLs (they'll open in your browser)
- Test API connectivity and responses

## Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| "Set SIMBRIEF_USER_ID" | Get your user ID from SimBrief account settings |
| "No flight plans found" | Create a flight plan on SimBrief website first |
| "Network error" | Check internet connection and SimBrief status |
| "Import errors" | Run `go mod tidy` to install dependencies |

### Debug Mode

Enable detailed logging to diagnose issues:

```bash
export SIMBRIEF_DEBUG=true
go run main.go
```

This will show:
- HTTP requests and responses
- API endpoints being called
- Response parsing details
- Error details and stack traces

### Getting Help

1. **Check the logs** with debug mode enabled
2. **Review the documentation** in the [docs/](../docs/) folder
3. **Examine the source code** in [pkg/](../pkg/) for implementation details
4. **Open an issue** on GitHub with debug output

## Next Steps

After exploring the examples:

1. **Read the Documentation**
   - [Main README](../README.md) - Complete SDK overview
   - [API Integration Guide](../docs/api-integration.md) - Best practices
   - [Flight Planning Guide](../docs/flight-planning.md) - Concepts and workflows

2. **Explore the Source Code**
   - [pkg/client/](../pkg/client/) - Client implementation
   - [pkg/types/](../pkg/types/) - Data structures and types

3. **Build Your Own Application**
   - Start with patterns from these examples
   - Add your specific requirements
   - Implement error handling and logging
   - Consider rate limiting and caching

4. **Contribute Back**
   - Found a bug? Open an issue
   - Have improvements? Submit a pull request
   - Need a feature? Let us know

## Example Progression

```
Basic Example
    ↓
Learn fundamentals
    ↓
Advanced Example  
    ↓
Understand production patterns
    ↓
Your Custom Application
```

Each example builds on concepts from the previous one, so we recommend starting with the basic example even if you're experienced with APIs.

## Performance Tips

- **Cache aircraft data** - Aircraft types don't change frequently
- **Implement retries** - Handle temporary network issues
- **Use appropriate timeouts** - Don't let requests hang indefinitely
- **Monitor rate limits** - Respect SimBrief's API usage policies
- **Log strategically** - Use debug mode for development, minimal logging for production

Happy flight planning! 🛫
