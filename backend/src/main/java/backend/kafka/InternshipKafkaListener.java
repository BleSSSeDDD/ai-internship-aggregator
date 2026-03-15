package backend.kafka;

import backend.service.InternshipService;
import com.aggregator.internship.CompanyInternship;
import com.google.protobuf.InvalidProtocolBufferException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
@Slf4j
public class InternshipKafkaListener {

    private final InternshipService internshipService;

    @KafkaListener(topics = "internships", groupId = "internship-db-consumer")
    public void listen(byte[] payload) {
        try {
            CompanyInternship internship = CompanyInternship.parseFrom(payload);

            if (internship == null || internship.getPositionName().isEmpty()) {
                log.warn("Received invalid internship data: {}", internship);
                throw new IllegalArgumentException("Invalid internship data");
            }

            internshipService.findOrCreateInternship(internship);

            log.debug("Successfully processed internship: {}", internship.getPositionName());

        } catch (InvalidProtocolBufferException e) {
            log.error("Failed to parse Protobuf message", e);
            throw new RuntimeException("Protobuf parse error", e);
        } catch (Exception e) {
            log.error("Failed to process internship message", e);
            throw e;
        }
    }
}