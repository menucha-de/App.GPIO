CFLAGS=-std=c99 -fPIC -MMD -MP
LDFLAGS=-shared
SOURCES=mica_gpio.c
TARGET=libmica-gpio.so
OBJS=$(SOURCES:.c=.o)

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(LDFLAGS) $(CFLAGS) -o $@ $(OBJS)

clean:
	rm -f $(OBJS) $(TARGET) *.d
