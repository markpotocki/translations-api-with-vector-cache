# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: translate.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    0,
    '',
    'translate.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0ftranslate.proto\x12\ttranslate\"T\n\x12TranslationRequest\x12\x0c\n\x04text\x18\x01 \x01(\t\x12\x17\n\x0fsource_language\x18\x02 \x01(\t\x12\x17\n\x0ftarget_language\x18\x03 \x01(\t\"*\n\x13TranslationResponse\x12\x13\n\x0btranslation\x18\x01 \x01(\t2X\n\nTranslator\x12J\n\tTranslate\x12\x1d.translate.TranslationRequest\x1a\x1e.translate.TranslationResponseb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'translate_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  DESCRIPTOR._loaded_options = None
  _globals['_TRANSLATIONREQUEST']._serialized_start=30
  _globals['_TRANSLATIONREQUEST']._serialized_end=114
  _globals['_TRANSLATIONRESPONSE']._serialized_start=116
  _globals['_TRANSLATIONRESPONSE']._serialized_end=158
  _globals['_TRANSLATOR']._serialized_start=160
  _globals['_TRANSLATOR']._serialized_end=248
# @@protoc_insertion_point(module_scope)
