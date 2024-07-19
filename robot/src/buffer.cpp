#include "buffer.h"

const uint8_t ACTION_ID_OFFSET = 0;
const uint8_t ACTION_ID_SIZE = 1;
const uint8_t X_JOINT_TRANSLATION_OFFSET = ACTION_ID_OFFSET + ACTION_ID_SIZE;
const uint8_t X_JOINT_TRANSLATION_SIZE = 4;
const uint8_t Y_JOINT_TRANSLATION_OFFSET = X_JOINT_TRANSLATION_OFFSET + X_JOINT_TRANSLATION_SIZE;
const uint8_t Y_JOINT_TRANSLATION_SIZE = 4;
const uint8_t Z_JOINT_TRANSLATION_OFFSET = Y_JOINT_TRANSLATION_OFFSET + Y_JOINT_TRANSLATION_SIZE;
const uint8_t Z_JOINT_TRANSLATION_SIZE = 4;
const uint8_t V_JOINT_TRANSLATION_OFFSET = Z_JOINT_TRANSLATION_OFFSET + Z_JOINT_TRANSLATION_SIZE;
const uint8_t V_JOINT_TRANSLATION_SIZE = 4;
const uint8_t W_JOINT_TRANSLATION_OFFSET = V_JOINT_TRANSLATION_OFFSET + V_JOINT_TRANSLATION_SIZE;
const uint8_t W_JOINT_TRANSLATION_SIZE = 4;

const uint8_t RESULT_CODE_SIZE = 1;

const uint8_t TOTAL_ACTION_WITH_TRANSLATIONS_SIZE = ACTION_ID_SIZE + X_JOINT_TRANSLATION_SIZE + Y_JOINT_TRANSLATION_SIZE + Z_JOINT_TRANSLATION_SIZE + V_JOINT_TRANSLATION_SIZE + W_JOINT_TRANSLATION_SIZE;
const uint8_t TOTAL_RESULT_WITH_TRANSLATIONS_SIZE = RESULT_CODE_SIZE + X_JOINT_TRANSLATION_SIZE + Y_JOINT_TRANSLATION_SIZE + Z_JOINT_TRANSLATION_SIZE + V_JOINT_TRANSLATION_SIZE + W_JOINT_TRANSLATION_SIZE;


void clear_buffer(uint8_t *buffer) {
  memset(buffer, 0, UART_BUFFER_SIZE);
}

RESULT_CODE load_translations_from_buffer(uint8_t *buffer, size_t buffer_len, JointTranslations *translations) {
    if (buffer_len < TOTAL_ACTION_WITH_TRANSLATIONS_SIZE) {
        return RESULT_INVALID_NUMBER_OF_PARAMETERS;
    }

    memcpy(translations->x.b, buffer+X_JOINT_TRANSLATION_OFFSET, X_JOINT_TRANSLATION_SIZE);
    memcpy(translations->y.b, buffer+Y_JOINT_TRANSLATION_OFFSET, Y_JOINT_TRANSLATION_SIZE);
    memcpy(translations->z.b, buffer+Z_JOINT_TRANSLATION_OFFSET, Z_JOINT_TRANSLATION_SIZE);
    memcpy(translations->v.b, buffer+V_JOINT_TRANSLATION_OFFSET, V_JOINT_TRANSLATION_SIZE);
    memcpy(translations->w.b, buffer+W_JOINT_TRANSLATION_OFFSET, W_JOINT_TRANSLATION_SIZE);

    return RESULT_OK;
}

size_t load_result_with_fallback_to_buffer(uint8_t *buffer, RESULT_CODE code, JointTranslations *translations) {
    buffer[0] = code;
    memcpy(buffer+X_JOINT_TRANSLATION_OFFSET, translations->x.b, X_JOINT_TRANSLATION_SIZE);
    memcpy(buffer+Y_JOINT_TRANSLATION_OFFSET, translations->y.b, Y_JOINT_TRANSLATION_SIZE);
    memcpy(buffer+Z_JOINT_TRANSLATION_OFFSET, translations->z.b, Z_JOINT_TRANSLATION_SIZE);
    memcpy(buffer+V_JOINT_TRANSLATION_OFFSET, translations->v.b, V_JOINT_TRANSLATION_SIZE);
    memcpy(buffer+W_JOINT_TRANSLATION_OFFSET, translations->w.b, W_JOINT_TRANSLATION_SIZE);

    return TOTAL_RESULT_WITH_TRANSLATIONS_SIZE;
}

size_t load_result_code_to_buffer(uint8_t *buffer, RESULT_CODE code) {
    buffer[0] = code;

    return RESULT_CODE_SIZE;
}