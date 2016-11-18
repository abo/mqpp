# mqpp
MQTT Packet Parser

Usage
---
```
    package main

    import "github.com/abo/mqpp"

    func main() {
        ...
        scanner := mqpp.NewScanner(r)
        for scanner.Scan() {
            p, err := scanner.Packet()
            if err != nil {
                
            }
            ....
        }

        if err := scanner.Err(); err != nil {
            
        }
    }

    
```