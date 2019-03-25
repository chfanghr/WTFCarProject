#include "ArduinoJson-v6.9.1.hpp"
#include <Arduino.h>
#include <Wire.h>

//#define DEBUG
using namespace ArduinoJson;
const int CommandGpio = 0, CommandIr = 1, CommandSerial = 2,
          CommandSerialHost = 3, CommandData = 4;
const int GpioPinmode = 0, GpioDigitalwrite = 1, GpioDigitalread = 2,
          GpioAnalogwrite = 3, GpioAnalogread = 4;
const int CommandSucceed = 1, CommandFailed = 0;
DynamicJsonDocument JsonMessage(100);
int CommandType;
int CommandMethod;
int CommandStatus = 1;
uint8_t PinNumber;
uint8_t PinValue;
int Param = 0;
const uint8_t i2c_addr = 9;

void setup() {
#ifdef DEBUG
  Serial.begin(9600);
#endif
  Wire.begin(i2c_addr);
  Wire.onReceive(receiveEvent);
}

void receiveEvent(int) {
  DeserializationError error = deserializeJson(JsonMessage, Wire);
  if (error) {
#ifdef DEBUG
    Serial.print(F("deserializeJson() failed: "));
    Serial.println(error.c_str());
#endif
    return;
  }

  CommandType = JsonMessage["type"];
  CommandMethod = JsonMessage["method"];
  PinNumber = JsonMessage["param"][0];
  PinValue = JsonMessage["param"][1];

  JsonMessage.clear();

  switch (CommandType) {
  case CommandGpio:
    // CommandGpio
    switch (CommandMethod) {
    case GpioPinmode:
      // GpioPinmode
      switch (PinValue) {
      case 0: // INPUT Mode
        pinMode(PinNumber, INPUT);
        break;

      case 1: // INPUT_PULLUP Mode
        pinMode(PinNumber, INPUT_PULLUP);
        break;
      case 2: // INPUT_PULLDOWN Mode
        pinMode(PinNumber, INPUT);
        break;
      case 3: // OUTPUT Mode
        pinMode(PinNumber, OUTPUT);
      }

    case GpioDigitalwrite:
      // GpioDigitalwrite
      switch (PinValue) {
      case HIGH:
        digitalWrite(PinNumber, HIGH);
        break;
      case LOW:
        digitalWrite(PinNumber, LOW);
        break;
      }
    case GpioDigitalread:
      // GpioDigitalread
      Param = digitalRead(PinNumber);
      break;
    case GpioAnalogwrite:
      // GpioAnalogwrite
      analogWrite(PinNumber, PinValue);
      break;
    case GpioAnalogread:
      // GpioAnalogread
      Param = analogRead(PinNumber);
      break;
    default:
      goto error;
    }
  case CommandIr:
    // CommandIr
  case CommandSerial:
    // CommandSerial
  case CommandSerialHost:
    // CommandSerialHost
  case CommandData:
    // CommandData
  default:
  error:
#ifdef DEBUG
    Serial.print("Operation failed.");
#endif
    CommandStatus = CommandFailed;
  }
  JsonMessage["type"] = CommandType;
  JsonMessage["status"] = CommandStatus;
  JsonArray CommandParam = JsonMessage.createNestedArray("param");
  CommandParam.add(Param);
  serializeJson(JsonMessage, Wire);
  CommandStatus = CommandSucceed;
}

void loop() {}
