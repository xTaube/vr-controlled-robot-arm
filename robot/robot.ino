#include <AccelStepper.h>
#include <MultiStepper.h>
#include <Servo.h>

#define DEG_PER_STEP 0.1125
#define STEP_PER_REVOLUTION (360 / DEG_PER_STEP)
#define MAX_SPEED 100
#define SPEED 100.0
#define Y_DIRECTION_PIN 2
#define Y_STEP_PIN 3
#define Z_DIRECTION_PIN 4
#define Z_STEP_PIN 5
#define X_DIRECTION_PIN 8
#define X_STEP_PIN 7
#define V_SERVO_PWM_PIN 9
#define W_SERVO_PWM_PIN 6
#define GRIPPER_MOTOR_B1_PIN 10
#define GRIPPER_MOTOR_B2_PIN 11

AccelStepper y_stepper(AccelStepper::DRIVER, Y_STEP_PIN, Y_DIRECTION_PIN);
AccelStepper z_stepper(AccelStepper::DRIVER, Z_STEP_PIN, Z_DIRECTION_PIN);
AccelStepper x_stepper(AccelStepper::DRIVER, X_STEP_PIN, X_DIRECTION_PIN);
Servo v_servo;
Servo w_servo;


void setup() {
  pinMode(GRIPPER_MOTOR_B1_PIN, OUTPUT);
  pinMode(GRIPPER_MOTOR_B2_PIN, OUTPUT);

  v_servo.write(90);
  v_servo.attach(V_SERVO_PWM_PIN);
  w_servo.write(90);
  w_servo.attach(W_SERVO_PWM_PIN);

  x_stepper.setCurrentPosition(0);
  x_stepper.setMaxSpeed(MAX_SPEED);
  x_stepper.setAcceleration(SPEED);

  y_stepper.setCurrentPosition(0);
  y_stepper.setMaxSpeed(MAX_SPEED);
  y_stepper.setAcceleration(SPEED);


  z_stepper.setCurrentPosition(0);
  z_stepper.setMaxSpeed(MAX_SPEED);
  z_stepper.setAcceleration(SPEED);
}


void loop() {}
