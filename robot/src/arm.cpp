#include "arm.h"

const float STEPS_PER_ONE_DEGREE = 1/DEG_PER_STEP;


RESULT_CODE move_arm(Arm *arm, JointTranslations *translations, JointTranslations *fallback) {
    long steps = translations->x.f*STEPS_PER_ONE_DEGREE;
    fallback->x.f = translations->x.f - steps*DEG_PER_STEP;
    arm->x_stepper->move(steps);

    steps = translations->y.f*STEPS_PER_ONE_DEGREE;
    fallback->y.f = translations->y.f - steps*DEG_PER_STEP;
    arm->y_stepper->move(steps);

    steps = translations->z.f*STEPS_PER_ONE_DEGREE;
    fallback->z.f = translations->z.f - steps*DEG_PER_STEP;
    arm->z_stepper->move(steps);

    fallback->v.f = (float)(int)translations->v.f - translations->v.f;
    arm->v_servo->write(arm->v_servo->read() + (int)translations->v.f);

    fallback->w.f = (float)(int)translations->w.f - translations->w.f;
    arm->w_servo->write(arm->w_servo->read() + (int)translations->w.f);

    return RESULT_OK;
}