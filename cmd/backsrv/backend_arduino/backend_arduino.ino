#include <Wire.h>
#include <Arduino.h>
#include "ArduinoJson-v6.9.1.hpp"

//#define DEBUG
using namespace ArduinoJson;

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
    Wire.onRequest(requestEvent);
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
        case 0:
            //CommandGpio
            switch (CommandMethod) {
                case 0:
                    //GpioPinmode
                    switch (PinValue) {
                        case 0:
                            pinMode(PinNumber, INPUT);
                            break;

                        case 1:
                            pinMode(PinNumber, INPUT_PULLUP);
                            break;
                        case 2:
                            pinMode(PinNumber, INPUT);
                            break;
                        case 3:
                            pinMode(PinNumber, OUTPUT);

                    }

                case 1:
                    //GpioDigitalwrite
                    switch (PinValue) {
                        case 0:
                            digitalWrite(PinNumber, HIGH);
                            break;
                        case 1:
                            digitalWrite(PinNumber, LOW);
                            break;
                    }
                case 2:
                    //GpioDigitalread
                    Param = digitalRead(PinNumber);
                    break;
                case 3:
                    //GpioAnalogwrite
                    analogWrite(PinNumber, PinValue);
                    break;
                case 4:
                    //GpioAnalogread
                    Param = analogRead(PinNumber);
                    break;
                default:
                    goto error;
            }
        case 1:
            //CommandIr
        case 2:
            //CommandSerial
        case 3:
            //CommandSerialHost
        case 4:
            //CommandData
        default:
        error:
            Serial.print("Operation failed.");
            CommandStatus = 0;
    }


}

void requestEvent() {
    JsonMessage["type"] = CommandType;
    JsonMessage["status"] = CommandStatus;
    JsonArray CommandParam = JsonMessage.createNestedArray("param");
    CommandParam.add(Param);
    serializeJson(JsonMessage, Wire);
}

void loop() {

}
