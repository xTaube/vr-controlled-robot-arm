#ifndef ARM_H
#define ARM_H

#include <AccelStepper.h>
#include <Servo.h>

#define DEG_PER_STEP 0.1125
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

typedef enum {
    MOVE_ACTION = 1
} ACTION_TYPE;

typedef enum {
    RESULT_OK = 1,
    RESULT_INVALID_NUMBER_OF_PARAMETERS = 2,
    RESULT_UNKNOWN_ACTION = 3
}  RESULT_CODE;

union JointTranslation {
    float f;
    byte b[4];
};

typedef struct {
    JointTranslation x;
    JointTranslation y;
    JointTranslation z;
    JointTranslation v;
    JointTranslation w;
} JointTranslations;


struct ArmState {
    bool is_calibrated;
};

struct Arm {
    AccelStepper *x_stepper;
    AccelStepper *y_stepper;
    AccelStepper *z_stepper;
    Servo *v_servo;
    Servo *w_servo;
    ArmState state;
};

RESULT_CODE move_arm(Arm *arm, JointTranslations *translations, JointTranslations *fallback);

#endif