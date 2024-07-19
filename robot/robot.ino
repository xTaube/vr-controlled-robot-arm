#include <AccelStepper.h>
#include <Servo.h>
#include <float.h>

#define BAUD_RATE 115200
#define DEG_PER_STEP 0.1125
#define UART_BUFFER_SIZE 128
#define MAX_SPEED 50
#define SPEED 50.0
#define Y_DIRECTION_PIN 2
#define Y_STEP_PIN 10
#define Z_DIRECTION_PIN 4
#define Z_STEP_PIN 5
#define X_DIRECTION_PIN 8
#define X_STEP_PIN 7
#define V_SERVO_PWM_PIN 9
#define W_SERVO_PWM_PIN 6
#define GRIPPER_MOTOR_B1_PIN 3
#define GRIPPER_MOTOR_B2_PIN 11
#define END_OF_TRANSMISSION 0x04

const float STEPS_PER_ONE_DEGREE = 1/DEG_PER_STEP;

AccelStepper y_stepper(AccelStepper::DRIVER, Y_STEP_PIN, Y_DIRECTION_PIN);
AccelStepper z_stepper(AccelStepper::DRIVER, Z_STEP_PIN, Z_DIRECTION_PIN);
AccelStepper x_stepper(AccelStepper::DRIVER, X_STEP_PIN, X_DIRECTION_PIN);
Servo v_servo;
Servo w_servo;

int x_pos = 1400;
int y_pos = 1000;
int z_pos = 1200;

uint8_t rbuffer[32];
uint8_t wbuffer[32];

typedef enum {
  MOVE_JOINTS_COMMAND = 1
} COMMAND_TYPE;

union Joint {
  float f;
  byte b[4];
};

typedef struct {
  Joint x;
  Joint y;
  Joint z;
  Joint v;
  Joint w;
} MoveJointsCommand;


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
}

void random_movement() {
  v_servo.write(random(5, 175));
  w_servo.write(random(5, 175));
  x_stepper.moveTo(x_pos);
  y_stepper.moveTo(y_pos);
  z_stepper.moveTo(z_pos);

  int running = 3;
  while (running > 0) {
      running = 0;
      if (x_stepper.distanceToGo() != 0) {
          running++;
          x_stepper.run();
      }

      if (y_stepper.distanceToGo() != 0) {
        running++;
        y_stepper.run();
      }

      if (z_stepper.distanceToGo() != 0) {
        running++;
        z_stepper.run();
      }
  }

  x_pos = random(-1000, 1000);
  y_pos = random(-1000, 1000);
  z_pos = random(-1200, 1200);
}

void load_to_move_joints_command(MoveJointsCommand *command) {
  memcpy(command->x.b, rbuffer+1, 4);
  memcpy(command->y.b, rbuffer+5, 4);
  memcpy(command->z.b, rbuffer+9, 4);
  memcpy(command->v.b, rbuffer+13, 4);
  memcpy(command->w.b, rbuffer+17, 4);
}


void send_back(MoveJointsCommand *command) {
  wbuffer[0] = 1;
  memcpy(wbuffer+1, command->x.b, 4);
  memcpy(wbuffer+5, command->y.b, 4);
  memcpy(wbuffer+9, command->z.b, 4);
  memcpy(wbuffer+13, command->v.b, 4);
  memcpy(wbuffer+17, command->w.b, 4);
  wbuffer[21] = 0x04; // Add EOT byte
  Serial.write(wbuffer, 22);
}

void loop() {
  if (Serial.available() > 0) {
    Serial.readBytesUntil(END_OF_TRANSMISSION, rbuffer, UART_BUFFER_SIZE);
    switch (rbuffer[0])
    {
      case MOVE_JOINTS_COMMAND:
          MoveJointsCommand *command = (MoveJointsCommand*) malloc(sizeof(MoveJointsCommand));
          load_to_move_joints_command(command);
          send_back(command);
          free(command);

        break;
    
      default:
        break;
    }
  }
}
