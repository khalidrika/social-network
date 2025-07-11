<template>
  <EventsList 
    :group-id="groupId"
    :can-create-event="canCreateEvent"
  />
</template>

<script setup>
import EventsList from './EventsList.vue'

defineProps({
  groupId: {
    type: String,
    required: true
  },
  canCreateEvent: {
    type: Boolean,
    default: false
  }
})
</script>


