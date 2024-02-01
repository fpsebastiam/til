Named tuples can be useful when customizing Enums. Here is an example:


```python
from collections import namedtuple
from enum import Enum


if __name__ == "__main__":
    TrafficLightTuple = namedtuple("TrafficLightTuple", "id, name, is_traffic_allowed")

    class TrafficLight(TrafficLightTuple, Enum):
        GREEN = TrafficLightTuple(1, "Green", True)
        YELLOW = TrafficLightTuple(2, "Yellow", True)
        RED = TrafficLightTuple(3, "Red", False)

        @classmethod
        def from_name(cls, name: str):
            for tl in cls:
                if name == tl.name:
                    return tl
            # handle non-matches here

    y = TrafficLight.YELLOW
    print(y.id, y.name, y.is_traffic_allowed)  # 2 Yellow true
    print(TrafficLight.from_name("Red"))  # TrafficLight.RED
```
I few observations:
- In the example above, I found the subclassing to be needed in some earlier versions of Python (e.g below 3.9)
- Enum members are singletons
- I experienced some performance hits when accesing enum properties or methods. The named tuple approach helped. But the Enum class provides a lot of extra functionality, so it might be worth for your use case.
- If the enum value does not matter, consider using auto()
- You can create aliases:
```python
from enum import Enum
class Shape(Enum):
    SQUARE = 2
    DIAMOND = 1
    CIRCLE = 3
    ALIAS_FOR_SQUARE = 2
```
