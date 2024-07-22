#ifndef BUFFER_H
#define BUFFER_H

#include "arm.h"

#define UART_BUFFER_SIZE 32

void clear_buffer(uint8_t *buffer);

RESULT_CODE load_translations_from_buffer(uint8_t *buffer, size_t buffer_len, JointTranslations *translations);

size_t load_result_with_fallback_to_buffer(uint8_t *buffer, RESULT_CODE code, JointTranslations *translations);

size_t load_result_code_to_buffer(uint8_t *buffer, RESULT_CODE code);

#endif