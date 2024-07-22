#include "src/arm.h"
#include "src/buffer.h"

#define BAUD_RATE 115200
#define END_OF_TRANSMISSION 0x04

AccelStepper y_stepper(AccelStepper::DRIVER, Y_STEP_PIN, Y_DIRECTION_PIN);
AccelStepper z_stepper(AccelStepper::DRIVER, Z_STEP_PIN, Z_DIRECTION_PIN);
AccelStepper x_stepper(AccelStepper::DRIVER, X_STEP_PIN, X_DIRECTION_PIN);
Servo v_servo;
Servo w_servo;

Arm arm = {
  &x_stepper,
  &y_stepper, 
  &z_stepper, 
  &v_servo, 
  &w_servo, 
  ArmState{false}
};

uint8_t buffer[UART_BUFFER_SIZE] = {0};
size_t loaded_bytes;
RESULT_CODE result_code;

void send_result(size_t size) {
  buffer[size] = END_OF_TRANSMISSION;
  Serial.write(buffer, size+1);
}

void setup() {
  // configure uart
  Serial.begin(BAUD_RATE, SERIAL_8E1);

  // configure gripper
  pinMode(GRIPPER_MOTOR_B1_PIN, OUTPUT);
  pinMode(GRIPPER_MOTOR_B2_PIN, OUTPUT);

  // configure v and w axes servo motors
  v_servo.attach(V_SERVO_PWM_PIN);
  w_servo.attach(W_SERVO_PWM_PIN);

  // configure x ax stepper motor
  x_stepper.setCurrentPosition(0);
  x_stepper.setMaxSpeed(MAX_SPEED);
  x_stepper.setAcceleration(SPEED);

  // configure y ax stepper motors
  y_stepper.setCurrentPosition(0);
  y_stepper.setMaxSpeed(MAX_SPEED);
  y_stepper.setAcceleration(SPEED);

  // configure z ax stepper motor
  z_stepper.setCurrentPosition(0);
  z_stepper.setMaxSpeed(MAX_SPEED);
  z_stepper.setAcceleration(SPEED);
  pinMode(13, OUTPUT);
}

void loop() {
  if (Serial.available() > 0) {
    loaded_bytes = Serial.readBytesUntil(END_OF_TRANSMISSION, buffer, UART_BUFFER_SIZE);
    switch (buffer[0])
    {
      case SET_NEW_ARM_POSITION:
        JointTranslations *translations = (JointTranslations*) malloc(sizeof(JointTranslations));
        result_code = load_translations_from_buffer(buffer, loaded_bytes, translations);
        clear_buffer(buffer);

        if (result_code == RESULT_INVALID_NUMBER_OF_PARAMETERS) {
          loaded_bytes = load_result_code_to_buffer(buffer, RESULT_INVALID_NUMBER_OF_PARAMETERS);
          send_result(loaded_bytes);
          break;
        }

        JointTranslations *fallback = (JointTranslations*) malloc(sizeof(JointTranslations));
        result_code = set_new_arm_position(&arm, translations, fallback);
        loaded_bytes = load_result_with_fallback_to_buffer(buffer, result_code, fallback);
        send_result(loaded_bytes);

        free(translations);
        free(fallback);
        break;
    
      default:
        loaded_bytes = load_result_code_to_buffer(buffer, RESULT_UNKNOWN_ACTION);
        send_result(loaded_bytes);
    }

    clear_buffer(buffer);
  }
  move_arm_steppers(&arm);
}
