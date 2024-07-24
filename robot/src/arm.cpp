#include "arm.h"


const float STEPS_PER_DEGREE = 1.0F/DEG_PER_STEP;
const float X_AX_GEAR_RATIO = 4.89F;
const float Y_AX_GEAR_RATIO = 6.0F;
const float Z_AX_GEAR_RATIO = 4.2F;

const float X_AX_STEPS_PER_DEGREE = STEPS_PER_DEGREE * X_AX_GEAR_RATIO;
const float Y_AX_STEPS_PER_DEGREE = STEPS_PER_DEGREE * Y_AX_GEAR_RATIO;
const float Z_AX_STEPS_PER_DEGREE = STEPS_PER_DEGREE * Z_AX_GEAR_RATIO;

const float X_AX_DEG_PER_STEP = DEG_PER_STEP/X_AX_GEAR_RATIO;
const float Y_AX_DEG_PER_STEP = DEG_PER_STEP/Y_AX_GEAR_RATIO;
const float Z_AX_DEG_PER_STEP = DEG_PER_STEP/Z_AX_GEAR_RATIO;


bool is_arm_in_move(Arm *arm) {
    return arm->x_stepper->isRunning() != 0 || arm->y_stepper->isRunning() || arm->z_stepper->isRunning();
}

void initialize_arm_motors(Arm *arm) {
    // configure gripper
  pinMode(GRIPPER_MOTOR_B1_PIN, OUTPUT);
  pinMode(GRIPPER_MOTOR_B2_PIN, OUTPUT);

  // configure v and w axes servo motors
  arm->v_servo->attach(V_SERVO_PWM_PIN);
  arm->w_servo->attach(W_SERVO_PWM_PIN);

  // configure x ax stepper motor
  arm->x_stepper->setCurrentPosition(0);
  arm->x_stepper->setMaxSpeed(MAX_SPEED);
  arm->x_stepper->setAcceleration(DEFAULT_SPEED);

  // configure y ax stepper motors
  arm->y_stepper->setCurrentPosition(0);
  arm->y_stepper->setMaxSpeed(MAX_SPEED);
  arm->y_stepper->setAcceleration(DEFAULT_SPEED);

  // configure z ax stepper motor
  arm->z_stepper->setCurrentPosition(0);
  arm->z_stepper->setMaxSpeed(MAX_SPEED);
  arm->z_stepper->setAcceleration(DEFAULT_SPEED);
}


RESULT_CODE set_new_arm_position(Arm *arm, JointsAngles *joints, JointsAngles *fallback) {
    if (!arm->state.is_calibrated && arm->state.mode != ARM_CALIBRATION_MODE) return RESULT_ARM_NOT_CALIBRATED;

    long steps = (long)round(joints->x*X_AX_STEPS_PER_DEGREE);
    fallback->x = (float)steps*X_AX_DEG_PER_STEP;
    arm->x_stepper->moveTo(steps);

    steps = (long)round(joints->y*Y_AX_STEPS_PER_DEGREE);
    fallback->y = (float)steps*Y_AX_DEG_PER_STEP;
    arm->y_stepper->moveTo(steps);

    steps = (long)round(joints->z*Z_AX_STEPS_PER_DEGREE);
    fallback->z = (float)steps*Z_AX_DEG_PER_STEP;
    arm->z_stepper->moveTo(steps);

    fallback->v = round(joints->v);
    arm->v_servo->write((int)fallback->v);

    fallback->w = round(joints->w);
    arm->w_servo->write((int)fallback->w);

    return RESULT_OK;
}

RESULT_CODE set_arm_speed(Arm *arm, float speed) {
    if (is_arm_in_move(arm)) return RESULT_ARM_IN_MOVE;
    if (speed > MAX_SPEED) return RESULT_BEYOND_MAX_SPEED_LIMIT;
    if (speed < MIN_SPEED) return RESULT_SPEED_TO_SLOW;

    arm->x_stepper->setAcceleration(speed);
    arm->y_stepper->setAcceleration(speed);
    arm->z_stepper->setAcceleration(speed);
    return RESULT_OK;
}

RESULT_CODE set_arm_current_position_as_reference(Arm *arm) {
    if (arm->state.mode != ARM_CALIBRATION_MODE) return RESULT_ARM_NOT_IN_CALIBRATION_MODE;
    if (is_arm_in_move(arm)) return RESULT_ARM_IN_MOVE;
    
    arm->x_stepper->setCurrentPosition(0);
    arm->y_stepper->setCurrentPosition(0);
    arm->z_stepper->setCurrentPosition(0);

    return RESULT_OK;
}

RESULT_CODE get_arm_current_position(Arm *arm, JointsAngles *position) {
    if (!arm->state.is_calibrated && arm->state.mode != ARM_CALIBRATION_MODE) return RESULT_ARM_NOT_CALIBRATED;

    position->x = (float)arm->x_stepper->currentPosition()*X_AX_DEG_PER_STEP;
    position->y = (float)arm->y_stepper->currentPosition()*Y_AX_DEG_PER_STEP;
    position->z = (float)arm->z_stepper->currentPosition()*Z_AX_DEG_PER_STEP;
    position->v = (float)arm->v_servo->read();
    position->w = (float)arm->v_servo->read();

    return RESULT_OK;
}

RESULT_CODE is_arm_calibrated(Arm *arm) {
    if (!arm->state.is_calibrated) return RESULT_ARM_NOT_CALIBRATED;
    return RESULT_OK;
}

void set_arm_mode(Arm *arm, ARM_MODE mode){
    arm->state.mode = mode;
}

void set_arm_calibration(Arm *arm, bool is_calibrated){
    arm->state.is_calibrated = is_calibrated;
}

void move_arm_steppers(Arm *arm) {
    if (arm->x_stepper->distanceToGo() != 0) arm->x_stepper->run();
    if (arm->y_stepper->distanceToGo() != 0) arm->y_stepper->run();
    if (arm->z_stepper->distanceToGo() != 0) arm->z_stepper->run();
}
