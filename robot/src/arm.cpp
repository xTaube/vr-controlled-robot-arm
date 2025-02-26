#include "arm.h"

#define X_AX_MIN_ANGLE -65
#define X_AX_MAX_ANGLE 120

#define Y_AX_MIN_ANGLE -180
#define Y_AX_MAX_ANGLE 5

#define V_AX_MIN_ANGLE -90
#define V_AX_MAX_ANGLE 90

#define W_AX_MIN_ANGLE -90
#define W_AX_MAX_ANGLE 90


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


bool Arm::is_in_move() {
    return this->x_stepper->isRunning() != 0 || this->y_stepper->isRunning() || this->z_stepper->isRunning();
}

void Arm::initialize_motors() {
    // configure gripper
  pinMode(GRIPPER_MOTOR_B1_PIN, OUTPUT);
  pinMode(GRIPPER_MOTOR_B2_PIN, OUTPUT);
  digitalWrite(GRIPPER_MOTOR_B1_PIN, LOW);
  digitalWrite(GRIPPER_MOTOR_B2_PIN, LOW);

  // configure v and w axes servo motors
  this->v_servo->attach(V_SERVO_PWM_PIN);
  this->w_servo->attach(W_SERVO_PWM_PIN);

  // configure x ax stepper motor
  this->x_stepper->setCurrentPosition(0);
  this->x_stepper->setMaxSpeed(MAX_SPEED);
  this->x_stepper->setAcceleration(DEFAULT_SPEED);

  // configure y ax stepper motors
  this->y_stepper->setCurrentPosition(-90*Y_AX_STEPS_PER_DEGREE);
  this->y_stepper->setMaxSpeed(MAX_SPEED);
  this->y_stepper->setAcceleration(DEFAULT_SPEED);

  // configure z ax stepper motor
  this->z_stepper->setCurrentPosition(0);
  this->z_stepper->setMaxSpeed(MAX_SPEED);
  this->z_stepper->setAcceleration(DEFAULT_SPEED);
}


RESULT_CODE Arm::set_new_position(JointsAngles *joints, JointsAngles *fallback) {
    if (!this->state.is_calibrated && this->state.mode != ARM_CALIBRATION_MODE) return RESULT_ARM_NOT_CALIBRATED;

    if (
        joints->x < X_AX_MIN_ANGLE || 
        joints->x > X_AX_MAX_ANGLE || 
        joints->y < Y_AX_MIN_ANGLE || 
        joints->y > Y_AX_MAX_ANGLE ||
        joints->w < W_AX_MIN_ANGLE ||
        joints->w > W_AX_MAX_ANGLE
    ) {
        return RESULT_ARM_INVALID_MOVE_RANGE;
    }

    long steps = (long)round(joints->x*X_AX_STEPS_PER_DEGREE);
    fallback->x = (float)steps*X_AX_DEG_PER_STEP;
    this->x_stepper->moveTo(steps);

    steps = (long)round(joints->y*Y_AX_STEPS_PER_DEGREE);
    fallback->y = (float)steps*Y_AX_DEG_PER_STEP;
    this->y_stepper->moveTo(steps);

    steps = (long)round(joints->z*Z_AX_STEPS_PER_DEGREE);
    fallback->z = (float)steps*Z_AX_DEG_PER_STEP;
    this->z_stepper->moveTo(steps);

    fallback->v = round(joints->v);
    int v = int(fallback->v) + 90;
    this->v_servo->write(v);

    fallback->w = round(joints->w);
    int w = 90 - int(fallback->w);
    if (w < 5) w = 5;
    this->w_servo->write(w);

    return RESULT_OK;
}

RESULT_CODE Arm::set_speed(float speed) {
    if (this->is_in_move()) return RESULT_ARM_IN_MOVE;
    if (speed > MAX_SPEED) return RESULT_BEYOND_MAX_SPEED_LIMIT;
    if (speed < MIN_SPEED) return RESULT_SPEED_TO_SLOW;

    this->x_stepper->setAcceleration(speed);
    this->y_stepper->setAcceleration(speed);
    this->z_stepper->setAcceleration(speed);
    return RESULT_OK;
}

RESULT_CODE Arm::set_current_position_as_reference() {
    if (this->state.mode != ARM_CALIBRATION_MODE) return RESULT_ARM_NOT_IN_CALIBRATION_MODE;
    if (this->is_in_move()) return RESULT_ARM_IN_MOVE;
    
    this->x_stepper->setCurrentPosition(0);
    this->y_stepper->setCurrentPosition(-90*Y_AX_STEPS_PER_DEGREE);
    this->z_stepper->setCurrentPosition(0);

    return RESULT_OK;
}

RESULT_CODE Arm::get_current_position(JointsAngles *position) {
    if (!this->state.is_calibrated && this->state.mode != ARM_CALIBRATION_MODE) return RESULT_ARM_NOT_CALIBRATED;

    position->x = (float)this->x_stepper->currentPosition()*X_AX_DEG_PER_STEP;
    position->y = (float)this->y_stepper->currentPosition()*Y_AX_DEG_PER_STEP;
    position->z = (float)this->z_stepper->currentPosition()*Z_AX_DEG_PER_STEP;
    position->v = (float)this->v_servo->read();
    position->w = (float)this->v_servo->read();

    return RESULT_OK;
}

RESULT_CODE Arm::is_calibrated() {
    if (!this->state.is_calibrated) return RESULT_ARM_NOT_CALIBRATED;
    return RESULT_OK;
}

void Arm::set_mode(ARM_MODE mode){
    this->state.mode = mode;
}

void Arm::set_calibration(bool is_calibrated){
    this->state.is_calibrated = is_calibrated;
}

void Arm::move_steppers() {
    if (this->x_stepper->distanceToGo() != 0) this->x_stepper->run();
    if (this->y_stepper->distanceToGo() != 0) this->y_stepper->run();
    if (this->z_stepper->distanceToGo() != 0) this->z_stepper->run();
}

void Arm::open_gripper() {
    analogWrite(GRIPPER_MOTOR_B1_PIN, 150);
    delay(100);

    analogWrite(GRIPPER_MOTOR_B1_PIN, 0);
}


void Arm::close_gripper() {
    analogWrite(GRIPPER_MOTOR_B2_PIN, 175);
    delay(100);

    analogWrite(GRIPPER_MOTOR_B2_PIN, 0);
}