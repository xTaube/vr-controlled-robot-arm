#ifndef ARM_H
#define ARM_H

#include <AccelStepper.h>
#include <Servo.h>

#define DEG_PER_STEP 0.1125
#define MAX_SPEED 1000.0
#define MIN_SPEED 50.0
#define DEFAULT_SPEED 50.0

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
    SET_NEW_ARM_POSITION = 1,
    SET_ARM_SPEED = 2,
    GET_ARM_CURRENT_POSITION = 3,
    CHECK_ARM_CALIBRATION = 4,
    START_CALIBRATION = 5,
    FINISH_CALIBRATION = 6,
    ABORT_CALIBRATION = 7,
    CHECK_ARM_IDLE = 8,
    OPEN_GRIPPER = 9,
    CLOSE_GRIPPER = 10
} ACTION_TYPE;

typedef enum {
    RESULT_OK = 1,
    RESULT_INVALID_NUMBER_OF_PARAMETERS = 10,
    RESULT_UNKNOWN_ACTION = 11,
    RESULT_ARM_NOT_CALIBRATED = 12,
    RESULT_BEYOND_MAX_SPEED_LIMIT = 13,
    RESULT_SPEED_TO_SLOW = 14,
    RESULT_ARM_IN_MOVE = 15,
    RESULT_ARM_NOT_IN_CALIBRATION_MODE = 16,
    RESULT_ARM_INVALID_MOVE_RANGE = 17
}  RESULT_CODE;

typedef enum {
    ARM_CALIBRATION_MODE,
    ARM_NORMAL_MODE,
} ARM_MODE;

typedef struct {
    float x;
    float y;
    float z;
    float v;
    float w;
} JointsAngles;


struct ArmState {
    bool is_calibrated;
    ARM_MODE mode;
};

struct Arm {
    AccelStepper *x_stepper;
    AccelStepper *y_stepper;
    AccelStepper *z_stepper;
    Servo *v_servo;
    Servo *w_servo;
    ArmState state;

    void initialize_motors();
    RESULT_CODE set_new_position(JointsAngles *translations, JointsAngles *fallback);
    RESULT_CODE set_speed(float speed);
    RESULT_CODE set_current_position_as_reference();
    RESULT_CODE get_current_position(JointsAngles *position);
    RESULT_CODE is_calibrated();
    bool is_in_move();
    void set_mode(ARM_MODE mode);
    void move_steppers();
    void set_calibration(bool is_calibrated);
    void open_gripper();
    void close_gripper();
};

#endif