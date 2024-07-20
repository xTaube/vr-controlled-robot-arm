#include "arm.h"


const float STEPS_PER_ONE_DEGREE = 1/DEG_PER_STEP;


RESULT_CODE set_new_arm_position(Arm *arm, JointTranslations *translations, JointTranslations *fallback) {
    long steps = (long)round(translations->x*STEPS_PER_ONE_DEGREE);
    fallback->x = (float)steps*DEG_PER_STEP - translations->x;
    arm->x_stepper->move(steps);

    steps = (long)round(translations->y*STEPS_PER_ONE_DEGREE);
    fallback->y = (float)steps*DEG_PER_STEP - translations->y;
    arm->y_stepper->move(steps);

    steps = (long)round(translations->z*STEPS_PER_ONE_DEGREE);
    fallback->z = (float)steps*DEG_PER_STEP - translations->z;
    arm->z_stepper->move(steps);

    fallback->v = round(translations->v) - translations->v;
    arm->v_servo->write(arm->v_servo->read() + (int)round(translations->v));

    fallback->w = round(translations->w) - translations->w;
    arm->w_servo->write(arm->w_servo->read() + (int)round(translations->w));

    return RESULT_OK;
}

void move_arm_steppers(Arm *arm) {
    if (arm->x_stepper->distanceToGo() != 0) arm->x_stepper->run();
    if (arm->y_stepper->distanceToGo() != 0) arm->y_stepper->run();
    if (arm->z_stepper->distanceToGo() != 0) arm->z_stepper->run();
}