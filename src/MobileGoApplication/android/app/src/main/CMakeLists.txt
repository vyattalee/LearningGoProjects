
add_library(lib_foo SHARED IMPORTED)
set_property(TARGET lib_foo PROPERTY IMPORTED_NO_SONAME 1)
set_target_properties(lib_foo PROPERTIES IMPORTED_LOCATION ${CMAKE_CURRENT_SOURCE_DIR}/../jniLibs/${CMAKE_ANDROID_ARCH_ABI}/libfoo.so)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../jniLibs/${CMAKE_ANDROID_ARCH_ABI}/)


target_link_libraries(native-lib lib_foo ${log-lib})