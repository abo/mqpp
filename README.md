# mqpp
MQTT Packet Parser

Usage
---
```
    package main

    import "github.com/abo/mqpp"

    func main() {
        ...
        splitter := mqpp.NewSplitter(r)
        for splitter.Scan() {
            p, err := splitter.Packet()
            if err != nil {
                
            }
            ....
        }

        if err := splitter.Err(); err != nil {
            
        }
    }

    
```